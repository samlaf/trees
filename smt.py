from bisect import bisect_right, bisect_left

hStarEmptyCache = ['0']

def HStarEmpty(n):
    if len(hStarEmptyCache) <= n:
        t = HStarEmpty(n-1)
        t = hash(t+t)
        assert len(hStarEmptyCache) == n
        hStarEmptyCache.append(t)
    return hStarEmptyCache[n]

def HStar2b(n, l, lo, hi, offset):
    t = hi - lo
    if n == 0:
        if t == 0:
            return '0'
        assert t == 1
        return '1'
    if t == 0:
        return HStarEmpty(n)
    split = (1 << (n-1)) + offset
    i = bisect_left(l, split, lo, hi)
    return hash(HStar2b(n-1,l,lo,i,offset) + HStar2b(n-1,l,i,hi,split))

def HStar2(n,l):
    l.sort()
    return HStar2b(n, l, 0, len(l), 0)

if __name__ == '__main__':
    # print(HStar2(3,[0,1]))
    # print(hStarEmptyCache)
    i = int(12345)
    i.bit_length